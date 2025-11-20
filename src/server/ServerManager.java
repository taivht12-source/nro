package server;

/*
 * @Author Coder: Nguyễn Tấn Tài
 * @Description: Ngọc Rồng Kiwi - Máy Chủ Chuẩn Teamobi 2025
 * @Group Zalo: https://zalo.me/g/toiyeuvietnam2025
 */
import services.shenron.Shenron_Manager;
import managers.ConsignShopManager;
import database.HistoryTransactionDAO;
import boss.BossManager.BossManager;
import boss.BossManager.OtherBossManager;
import boss.BossManager.TreasureUnderSeaManager;
import boss.BossManager.SnakeWayManager;
import boss.BossManager.RedRibbonHQManager;
import boss.BossManager.GasDestroyManager;
import boss.BossManager.YardartManager;
import boss.BossManager.SkillSummonedManager;
import interfaces.ISession;
import network.Network;
import network.MyKeyHandler;
import network.MySession;
import services.dungeon.NgocRongNamecService;
import utils.Logger;
import utils.TimeUtil;

import java.util.*;
import matches.tournament.The23rdMartialArtCongressManager;
import matches.tournament.DeathOrAliveArenaManager;
import matches.tournament.WorldMartialArtsTournamentManager;
import network.MessageSendCollect;
import managers.SuperRankManager;
import interfaces.ISessionAcceptHandler;
import boss.BossManager.BrolyManager;
import boss.BossManager.FinalBossManager;
import javax.swing.JFrame;
import javax.swing.JPanel;
import player.Service.ClanService;

public class ServerManager {

    public static String timeStart;

    public static final Map<Object, Object> CLIENTS = new HashMap<>();
    public static String NAME_SERVER = "Ngọc Rồng Black Goku"; // Tên Máy Chủ
    public static String DOMAIN = "https://nroblackgoku.click/"; // Domain Truy Cập
    public static String NAME = "NRO Black Goku"; // Name Khi Vào Giao Diện Game
    public static String IP = "NgocRongOnline"; // IPs - Không Cần Sửa
    public static int PORT = 14445; // Port
    public static int EVENT_SEVER = 0;

    private static ServerManager instance;

    public static boolean isRunning;

    public void init() {
        Manager.gI();
        HistoryTransactionDAO.deleteHistory();
    }

    public static ServerManager gI() {
        if (instance == null) {
            instance = new ServerManager();
            instance.init();
        }
        return instance;
    }

    public static void main(String[] args) {
        timeStart = TimeUtil.getTimeNow("dd/MM/yyyy HH:mm:ss");
        ServerManager.gI().run();
    }

   public void run() {
    isRunning = true;
    activeServerSocket();
    new Thread(NgocRongNamecService.gI(), "Update NRNM").start();
    new Thread(SuperRankManager.gI(), "Update Super Rank").start();
    new Thread(The23rdMartialArtCongressManager.gI(), "Update DHVT23").start();
    new Thread(DeathOrAliveArenaManager.gI(), "Update Võ Đài Sinh Tử").start();
    new Thread(WorldMartialArtsTournamentManager.gI(), "Update WMAT").start();
    new Thread(Shenron_Manager.gI(), "Update Shenron").start();
    BossManager.gI().loadBoss();
    Manager.MAPS.forEach(map.Map::initBoss);
    List<Runnable> PhoBan = new ArrayList<>();
    PhoBan.add(BossManager.gI());
    PhoBan.add(YardartManager.gI());
    PhoBan.add(FinalBossManager.gI());
    PhoBan.add(SkillSummonedManager.gI());
    PhoBan.add(BrolyManager.gI());
    PhoBan.add(OtherBossManager.gI());
    PhoBan.add(RedRibbonHQManager.gI());
    PhoBan.add(TreasureUnderSeaManager.gI());
    PhoBan.add(SnakeWayManager.gI());
    PhoBan.add(GasDestroyManager.gI());
    PhoBan.forEach(PhoBantarget -> new Thread(PhoBantarget, "Update PB").start());
    }

    public void activeServerSocket() {
        try {
            Network.gI().init().setAcceptHandler(new ISessionAcceptHandler() {
                @Override
                public void sessionInit(ISession is) {
                    if (!canConnectWithIp(is.getIP())) {
                        is.disconnect();
                        return;
                    }
                    is.setMessageHandler(Controller.gI())
                            .setSendCollect(new MessageSendCollect())
                            .setKeyHandler(new MyKeyHandler())
                            .startCollect().startQueueHandler();
                }

                @Override
                public void sessionDisconnect(ISession session) {
                    Client.gI().kickSession((MySession) session);
                    disconnect((MySession) session);
                }
            }).setTypeSessionClone(MySession.class)
              .setDoSomeThingWhenClose(() -> {
                  Logger.error("SERVER CLOSE\n");
                  System.exit(0);
              })
              .start(PORT);
        } catch (Exception e) {
            Logger.error("Lỗi khi khởi động máy chủ: " + e.getMessage());
        }
    }
    private boolean canConnectWithIp(String ipAddress) {
        Object o = CLIENTS.get(ipAddress);
        if (o == null) {
            CLIENTS.put(ipAddress, 1);
            return true;
        } else {
            int n = Integer.parseInt(String.valueOf(o));
            if (n < Manager.MAX_PER_IP) {
                n++;
                CLIENTS.put(ipAddress, n);
                return true;
            } else {
                return false;
            }
        }
    }
    public void disconnect(MySession session) {
        Object o = CLIENTS.get(session.getIP());
        if (o != null) {
            int n = Integer.parseInt(String.valueOf(o));
            n--;
            if (n < 0) {
                n = 0;
            }
            CLIENTS.put(session.getIP(), n);
        }
    }
    public void close() {
        isRunning = false;
        try {
            ClanService.gI().close();
        } catch (Exception e) {
            Logger.error("Lỗi save clan!\n");
        }
        try {
            ConsignShopManager.gI().save();
        } catch (Exception e) {
            Logger.error("Lỗi save shop ký gửi!\n");
        }
        Client.gI().close();
        Logger.success("SUCCESSFULLY MAINTENANCE!\n");
        System.exit(0);
    }
}
